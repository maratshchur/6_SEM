import os
import time
import json
import sys
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives.hashes import SHA256
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import base64
import secrets

def generate_key():
    return secrets.token_bytes(16)


def encrypt(key, plaintext):
    iv = os.urandom(12)
    cipher = Cipher(algorithms.AES(key), modes.GCM(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    ciphertext = encryptor.update(plaintext.encode()) + encryptor.finalize()
    return iv + encryptor.tag + ciphertext

def decrypt(key, ciphertext):
    iv = ciphertext[:12]
    tag = ciphertext[12:28]
    ciphertext = ciphertext[28:]
    cipher = Cipher(algorithms.AES(key), modes.GCM(iv, tag), backend=default_backend())
    decryptor = cipher.decryptor()
    return decryptor.update(ciphertext) + decryptor.finalize()


def derive_key_from_password(password, salt):
    kdf = PBKDF2HMAC(
        algorithm=SHA256(),
        length=16,
        salt=salt,
        iterations=100000,
        backend=default_backend()
    )
    return kdf.derive(password.encode())


def hash_password(password, salt):
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA512(),
        length=32,
        salt=salt,
        iterations=200000,
        backend=default_backend()
    )
    return base64.b64encode(kdf.derive(password.encode())).decode()

def verify_password(password, salt, hashed_password):
    expected_hash = hash_password(password, salt)
    return hashed_password == expected_hash


KDC_SECRET_KEY = generate_key()

KDC_DATABASE = {
    "users": {
        "user": {
            "salt": os.urandom(16),
            "hashed_password": None,
            "key": None,
        }
    },
    "services": {
        "server": generate_key(),
    }
}

def generate_password(length=12):
    alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
    return ''.join(secrets.choice(alphabet) for _ in range(length))

for username, user_data in KDC_DATABASE["users"].items():
    salt = user_data["salt"]
    password = generate_password()
    print(f"Generated password for {username}: {password}")
    user_data["hashed_password"] = hash_password(password, salt)
    user_data["key"] = derive_key_from_password(password, salt)


class KDC:
    def __init__(self, db, secret_key):
        self.db = db
        self.secret_key = secret_key  

    def authenticate_user(self, username, password):
        if username not in self.db["users"]:
            print("Error: User not found in KDC database!")
            sys.exit(1)

        user_data = self.db["users"][username]
        if not verify_password(password, user_data["salt"], user_data["hashed_password"]):
            print("Error: Invalid password!")
            sys.exit(1)

        session_key = generate_key()

        tgt = {
            "username": username,
            "session_key": session_key.hex(),
            "timestamp": int(time.time()),
            "expires_in": 300
        }

        tgt_encrypted = encrypt(self.secret_key, json.dumps(tgt))
        return tgt_encrypted, session_key

    def get_service_ticket(self, tgt_encrypted, service_name):
        try:
            tgt = json.loads(decrypt(self.secret_key, tgt_encrypted).decode())
        except Exception:
            print("Error: Failed to decode TGT!")
            sys.exit(1)

        if int(time.time()) > tgt["timestamp"] + tgt["expires_in"]:
            print("Error: TGT has expired!")
            sys.exit(1)

        session_key = bytes.fromhex(tgt["session_key"])

        if service_name not in self.db["services"]:
            print("Error: Requested service not found!")
            sys.exit(1)
        
        service_key = self.db["services"][service_name]

        service_ticket = {
            "username": tgt["username"],
            "session_key": session_key.hex(),
            "timestamp": int(time.time()),
            "expires_in": 300
        }

        service_ticket_encrypted = encrypt(service_key, json.dumps(service_ticket))
        return service_ticket_encrypted, session_key

class Service:
    def __init__(self, service_name, key):
        self.service_name = service_name
        self.key = key

    def validate_ticket(self, ticket_encrypted):
        try:
            ticket = json.loads(decrypt(self.key, ticket_encrypted).decode())
        except Exception:
            print("Error: Failed to decrypt service ticket!")
            sys.exit(1)

        if int(time.time()) > ticket["timestamp"] + ticket["expires_in"]:
            print("Error: Service ticket has expired!")
            sys.exit(1)

        print(f"[Service] Validated ticket: {ticket}")
        return ticket


class Client:
    def __init__(self, username, password, kdc):
        self.username = username
        self.password = password
        self.kdc = kdc
        self.session_key = None

    def authenticate(self):
        print(f"[Client] Authenticating as {self.username}")
        try:
            self.tgt, self.session_key = self.kdc.authenticate_user(self.username, self.password)
            print(f"[Client] Received TGT: {self.tgt}")
        except SystemExit:
            print("[Client] Authentication terminated due to an error.")
            raise

    def request_service(self, service_name):
        print(f"[Client] Requesting access to {service_name}")
        try:
            service_ticket, service_session_key = self.kdc.get_service_ticket(self.tgt, service_name)
            print(f"[Client] Received service ticket: {service_ticket}")
            return service_ticket
        except SystemExit:
            print("[Client] Service request terminated due to an error.")
            raise


if __name__ == "__main__":
    kdc = KDC(KDC_DATABASE, KDC_SECRET_KEY)
    server = Service("server", KDC_DATABASE["services"]["server"])

    user = Client("user", password, kdc)

    try:
        user.authenticate()
    except SystemExit:
        print("[Main] Program terminated during authentication.")
        sys.exit(1)

    try:
        service_ticket = user.request_service("server")
    except SystemExit:
        print("[Main] Program terminated during service request.")
        sys.exit(1)

    if service_ticket:
        try:
            server.validate_ticket(service_ticket)
        except SystemExit:
            print("[Main] Program terminated during ticket validation.")
            sys.exit(1)