import redis

import settings


redisdb = redis.StrictRedis(host=settings.REDIS_HOST, port=settings.REDIS_PORT,
                            db=settings.REDIS_DB_COUNT)
