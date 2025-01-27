def vigenere_cipher(text, key, encrypt=True):
    result = ''
    key_len = len(key)
    key_index = 0
    for char in text:
        if char.isalpha():
            start = ord('a') if char.islower() else ord('A')
            key_char = key[key_index % key_len].lower()
            shift = ord(key_char) - ord('a')
            shifted_char = chr((ord(char) - start + (shift if encrypt else -shift)) % 26 + start)
            key_index += 1
        else:
            shifted_char = char
        result += shifted_char
    return result