from pathlib import Path
import os



BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = 'c^z+a@#25hwp*tos+#oe^$c)9$q!!)*wg0$82!n=orc&@o_ql_'

DEBUG = True

ALLOWED_HOSTS = []

INSTALLED_APPS = ["db_app"]

MIDDLEWARE = []

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
        'NAME': os.path.join(BASE_DIR, "..","..","sqlite.db"),
    }
}


LANGUAGE_CODE = 'en-us'

TIME_ZONE = 'UTC'

USE_I18N = True

USE_L10N = True

USE_TZ = True
