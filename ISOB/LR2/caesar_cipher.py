def caesar_cipher(text, shift, encrypt=True):
    result = ''
    for char in text:
        if char.isalpha():
            if 'а' <= char <= 'я' or 'А' <= char <= 'Я':
                start = ord('а') if char.islower() else ord('А')
                alphabet_size = 32
            else:
                start = ord('a') if char.islower() else ord('A')
                alphabet_size = 26 
            
            shifted_char = chr((ord(char) - start + (shift if encrypt else -shift)) % alphabet_size + start)
        elif char.isdigit():
            shifted_char = str((int(char) + (shift if encrypt else -shift)) % 10)
        else:
            shifted_char = char
        result += shifted_char
    return result
