from redis import StrictRedis
from config.config import PLUGIN_CONFIG

import logging


class RedisClient():

    def __init__(self, **kwargs):

        self.host = None
        self.port = None
        if ":" in PLUGIN_CONFIG["RedisAddress"]:
            self.host, self.port = PLUGIN_CONFIG["RedisAddress"].split(":")
        self.connection = None  # database
        self._db_connection()

    def __del__(self):
        self.connection = None

    def _db_connection(self):
        try:
            if self.host is not None and self.port is not None:
                self.connection = StrictRedis(
                    host=self.host, port=self.port, db=PLUGIN_CONFIG["Db"], socket_timeout=PLUGIN_CONFIG["SocketTimeout"])
            else:
                # default host = "localhost", port = 6379
                self.connection = StrictRedis(
                    db=PLUGIN_CONFIG["Db"], socket_timeout=PLUGIN_CONFIG["SocketTimeout"])
        except Exception as err:
            logging.error(
                "Unable to connect to redis database. Error:{e}".format(e=err))

    def execute_command(self, *args, **kwargs):
        try:
            if self.connection:
                value = self.connection.execute_command(*args, **kwargs)
                if value is not None:
                    return value
        except Exception as err:
            logging.error("Unable to execute command. options: [{arg}, {kwarg}]. Error:{e}".format(
                arg=args, kwarg=kwargs, e=err))

    def set(self, key, value):
        try:
            if self.connection:
                self.connection.set(key, value)
        except Exception as err:
            logging.error(
                "unable to set key {key} in redis db. Error:{e}".format(key=key, e=err))

    def get(self, key):
        try:
            if self.connection:
                value = self.connection.get(key)
                if value:
                    return value.decode('utf-8')
        except Exception as err:
            logging.error(
                "unable to get key {key} in redis db. Error:{e}".format(key=key, e=err))

    def keys(self, pattern=None):
        decoded_keys = []
        try:
            if self.connection:
                keys = self.connection.keys(pattern)
                for value in keys:
                    decoded_keys.append(value.decode('utf-8'))
        except Exception as err:
            logging.error("unable to get redis keys for pattern {pat} in redis db. Error:{e}".format(
                pat=pattern, e=err))
        finally:
            return decoded_keys

    def delete(self, *keys):
        try:
            if self.connection:
                self.connection.delete(*keys)
        except Exception as err:
            logging.error(
                "unable to delete keys {keys} in redis db. Error:{e}".format(key=keys, e=err))

    def exists(self, *keys):
        exists = 0
        try:
            if self.connection:
                exists = self.connection.exists(*keys)
        except Exception as err:
            logging.error(
                "unable to exists keys {keys} in redis db. Error:{e}".format(key=keys, e=err))
        finally:
            return exists
