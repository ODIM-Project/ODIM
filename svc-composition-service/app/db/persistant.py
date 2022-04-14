#(C) Copyright [2022] American Megatrends International LLC
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

from redis import StrictRedis
from config.config import CONFIG_DATA

import logging


# Default Redis-on-disk
class RedisClient():
    def __init__(self, redis_address=None):

        self.host = None
        self.port = None
        if redis_address and ":" in redis_address:
            self.host, self.port = redis_address.split(":")
        elif ":" in CONFIG_DATA["RedisOnDiskAddress"]:
            self.host, self.port = CONFIG_DATA["RedisOnDiskAddress"].split(
                ":")
        self.connection = None  # database
        self._db_connection()

    def __del__(self):
        self.connection = None

    def _db_connection(self):
        try:
            if self.host is not None and self.port is not None:
                self.connection = StrictRedis(
                    host=self.host,
                    port=self.port,
                    db=CONFIG_DATA["Db"],
                    socket_timeout=CONFIG_DATA["SocketTimeout"])
            else:
                # default host = "localhost", port = 6379
                self.connection = StrictRedis(
                    db=CONFIG_DATA["Db"],
                    socket_timeout=CONFIG_DATA["SocketTimeout"])
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
            logging.error(
                "Unable to execute command. options: [{arg}, {kwarg}]. Error:{e}"
                .format(arg=args, kwarg=kwargs, e=err))

    def set(self, key, value):
        try:
            if self.connection:
                self.connection.set(key, value)
        except Exception as err:
            logging.error(
                "unable to set key {key} in redis db. Error:{e}".format(
                    key=key, e=err))

    def sadd(self, key, value):
        try:
            if self.connection:
                self.connection.sadd(key, value)
        except Exception as err:
            logging.error(
                "unable to sadd key {key} in redis db. Error:{e}".format(
                    key=key, e=err))

    def get(self, key):
        try:
            if self.connection:
                value = self.connection.get(key)
                if value:
                    return value.decode('utf-8')
        except Exception as err:
            logging.error(
                "unable to get key {key} in redis db. Error:{e}".format(
                    key=key, e=err))

    def smembers(self, key):
        decoded_value = []
        try:
            if self.connection:
                value = self.connection.smembers(key)
                for mem in value:
                    decoded_value.append(mem.decode('utf-8'))
        except Exception as err:
            logging.error(
                "unable to smembers key {key} in redis db. Error:{e}".format(
                    key=key, e=err))
        return decoded_value

    def keys(self, pattern=None):
        decoded_keys = []
        try:
            if self.connection:
                keys = self.connection.keys(pattern)
                for value in keys:
                    decoded_keys.append(value.decode('utf-8'))
        except Exception as err:
            logging.error(
                "unable to get redis keys for pattern {pat} in redis db. Error:{e}"
                .format(pat=pattern, e=err))
        return decoded_keys

    def delete(self, *keys):
        try:
            if self.connection:
                self.connection.delete(*keys)
        except Exception as err:
            logging.error(
                "unable to delete keys {key} in redis db. Error:{e}".format(
                    key=keys, e=err))

    def srem(self, key, *value):
        try:
            if self.connection:
                self.connection.srem(key, *value)
        except Exception as err:
            logging.error(
                "unable to delete keys {key} in redis db. Error:{e}".format(
                    key=key, e=err))

    def exists(self, *keys):
        exists = 0
        try:
            if self.connection:
                exists = self.connection.exists(*keys)
        except Exception as err:
            logging.error(
                "unable to exists keys {key} in redis db. Error:{e}".format(
                    key=keys, e=err))
        return exists

    def pipeline(self):
        pipe = None
        try:
            if self.connection:
                pipe = self.connection.pipeline()
        except Exception as err:
            logging.error(
                "Unable to create a pipeline in redis db. Error: {e}".format(
                    e=err))
        return pipe


RedisDb = RedisClient