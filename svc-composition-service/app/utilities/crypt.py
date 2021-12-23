from Crypto.Hash import SHA512
from Crypto.Cipher import PKCS1_OAEP
from Crypto.PublicKey import RSA
import base64
import os
import logging


class Crypt():
    def __init__(self, public_key_path, private_key_path):
        self.public_key = None
        self.private_key = None
        self.hash_object = SHA512.new()
        try:
            if os.path.exists(public_key_path):
                self.public_key = RSA.importKey(open(public_key_path).read())
            if os.path.exists(private_key_path):
                self.private_key = RSA.importKey(open(private_key_path).read())
        except Exception as err:
            logging.error(
                "Unable to parse certificates. Error: {e}".format(e=err))

    def encrypt(self, message):
        encode_text = ""
        try:
            if self.public_key:
                cipher = PKCS1_OAEP.new(self.public_key,
                                        hashAlgo=self.hash_object)
                encrypted_text = cipher.encrypt(message)
                encode_text = base64.b64encode(encrypted_text)
        except Exception as err:
            logging.error("Unable to encrypt text. Error: {e}".format(e=err))
        finally:
            return encode_text

    def decrypt(self, decrypt_text):
        result = ""
        try:
            if self.private_key:
                cipher = PKCS1_OAEP.new(self.private_key,
                                        hashAlgo=self.hash_object)
                decode_text = decode_text = base64.b64decode(decrypt_text)

                result = cipher.decrypt(decode_text)

                if isinstance(result, bytes):
                    result = result.decode('utf-8')
        except Exception as err:
            logging.error("Unable to Decrypt text. Error: {e}".format(e=err))
        finally:
            return result
