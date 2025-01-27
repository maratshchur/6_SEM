from caesar_cipher import caesar_cipher
from vigenere_cipher import vigenere_cipher

def main():
    with open("input.txt", "r") as file:
        input_text = file.read()

    cipher_choice = input("Выберите метод шифрования (1 - Цезарь, 2 - Виженер): ")

    if cipher_choice == "1":
        key = int(input("Введите ключ для шифра Цезаря (целое число): "))
        ciphered_text = caesar_cipher(input_text, key)
        ciphered_filename = "caesar_ciphered.txt"
    elif cipher_choice == "2":
        key = input("Введите ключ для шифра Виженера (строка): ")
        ciphered_text = vigenere_cipher(input_text, key)
        ciphered_filename = "vigenere_ciphered.txt"
    else:
        print("Некорректный выбор метода шифрования.")
        return

    with open(ciphered_filename, "w") as file:
        file.write(ciphered_text)
    print(f"Зашифрованный текст записан в файл: {ciphered_filename}")


    if cipher_choice == "1":
        deciphered_text = caesar_cipher(ciphered_text, key, encrypt=False)
        deciphered_filename = "caesar_deciphered.txt"
    elif cipher_choice == "2":
        deciphered_text = vigenere_cipher(ciphered_text, key, encrypt=False)
        deciphered_filename = "vigenere_deciphered.txt"

    with open(deciphered_filename, "w") as file:
        file.write(deciphered_text)
    print(f"Расшифрованный текст записан в файл: {deciphered_filename}")



if __name__ == "__main__":
    main()