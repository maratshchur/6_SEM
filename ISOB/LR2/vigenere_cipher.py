def vigenere_cipher(text, key, encrypt=True):
    result = ''
    key_len = len(key)
    key_index = 0
    
    for char in text:
        if char.isalpha():
            if 'а' <= char <= 'я' or 'А' <= char <= 'Я':  
                start = ord('а') if char.islower() else ord('А')
                alphabet_size = 32
            else:  
                start = ord('a') if char.islower() else ord('A')
                alphabet_size = 26
            key_char = key[key_index % key_len].lower()
            if 'а' <= key_char <= 'я':
                shift = ord(key_char) - ord('а')
            else: 
                shift = ord(key_char) - ord('a')
            
            shifted_char = chr((ord(char) - start + (shift if encrypt else -shift)) % alphabet_size + start)
            key_index += 1
        else:
            shifted_char = char
        result += shifted_char
    return result
