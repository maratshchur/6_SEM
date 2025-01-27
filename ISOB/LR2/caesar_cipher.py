def caesar_cipher(text, shift, encrypt=True):
    result = ''
    for char in text:
        if char.isalpha():
            start = ord('a') if char.islower() else ord('A')
            shifted_char = chr((ord(char) - start + (shift if encrypt else -shift)) % 26 + start)
        elif char.isdigit():
            shifted_char = str((int(char) + (shift if encrypt else -shift)) % 10)
        else:
            shifted_char = char
        result += shifted_char
    return result

