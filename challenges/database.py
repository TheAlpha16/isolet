from sqlalchemy import create_engine, Column, Integer, String, Text, Boolean, ForeignKey, ARRAY, TIMESTAMP, BigInteger
from sqlalchemy.dialects.postgresql import ENUM, INET
from sqlalchemy.orm import declarative_base, sessionmaker
import time
import os

Base = declarative_base()
chall_type_enum = ENUM('static', 'dynamic', 'on-demand', name='chall_type', create_type=True)

class Category(Base):
    __tablename__ = 'categories'
    
    category_id = Column(Integer, primary_key=True, autoincrement=True)
    category_name = Column(Text, nullable=False, unique=True)

class User(Base):
    __tablename__ = 'users'
    
    userid = Column(BigInteger, primary_key=True, autoincrement=True)
    email = Column(Text, nullable=False, unique=True)
    username = Column(Text, nullable=False, unique=True)
    score = Column(Integer, default=0)
    rank = Column(Integer, default=3)
    password = Column(String(100), nullable=False)
    teamid = Column(BigInteger, default=-1)

class Team(Base):
    __tablename__ = 'teams'
    
    teamid = Column(BigInteger, primary_key=True, autoincrement=True)
    teamname = Column(Text, nullable=False, unique=True)
    captain = Column(BigInteger, ForeignKey('users.userid'), nullable=False)
    members = Column(ARRAY(BigInteger), default=[])
    password = Column(String(100), nullable=False)

class Flag(Base):
    __tablename__ = 'flags'
    
    flagid = Column(BigInteger, primary_key=True, autoincrement=True)
    userid = Column(BigInteger, ForeignKey('users.userid'), nullable=False)
    level = Column(Integer)
    password = Column(Text)
    flag = Column(Text, nullable=False)
    port = Column(Integer, nullable=False)
    verified = Column(Boolean, default=False)
    hostname = Column(Text)
    deadline = Column(Integer, default=2526249600)
    extended = Column(Integer, default=1)

class Sublog(Base):
    __tablename__ = 'sublogs'
    
    sid = Column(BigInteger, primary_key=True, autoincrement=True)
    chall_id = Column(Integer, ForeignKey('challenges.chall_id'), nullable=False)
    userid = Column(BigInteger, ForeignKey('users.userid'), nullable=False)
    teamid = Column(BigInteger, ForeignKey('teams.teamid'), nullable=False)
    flag = Column(Text, nullable=False)
    correct = Column(Boolean, nullable=False)
    ip = Column(INET, nullable=False)
    subtime = Column(TIMESTAMP, default='NOW()')

class Running(Base):
    __tablename__ = 'running'
    
    runid = Column(BigInteger, primary_key=True, autoincrement=True)
    userid = Column(BigInteger, ForeignKey('users.userid'), nullable=False)
    level = Column(Integer)

class ToVerify(Base):
    __tablename__ = 'toverify'
    
    vid = Column(BigInteger, primary_key=True, autoincrement=True)
    email = Column(Text, nullable=False, unique=True)
    username = Column(Text, nullable=False, unique=True)
    password = Column(String(100), nullable=False)
    timestamp = Column(TIMESTAMP, default='NOW()')

class Challenge(Base):
    __tablename__ = 'challenges'
    
    chall_id = Column(Integer, primary_key=True, autoincrement=True)
    chall_name = Column(Text, nullable=False)
    category_id = Column(Integer, ForeignKey('categories.category_id'), nullable=False)
    prompt = Column(Text)
    flag = Column(Text)
    type = Column(chall_type_enum, nullable=False, default='static')
    points = Column(Integer, default=100)
    files = Column(ARRAY(Text), default=[])
    hints = Column(ARRAY(BigInteger), default=[])
    solves = Column(Integer, default=0)
    author = Column(Text, default='anonymous')
    visible = Column(Boolean, default=False)
    tags = Column(ARRAY(Text))
    port = Column(Integer, default=0)
    subd = Column(Text, default='localhost')
    cpu = Column(Integer, default=5)
    mem = Column(Integer, default=10)

class Hint(Base):
    __tablename__ = 'hints'
    
    hid = Column(BigInteger, primary_key=True, autoincrement=True)
    chall_id = Column(Integer, ForeignKey('challenges.chall_id'), nullable=False)
    hint = Column(Text, nullable=False)
    cost = Column(Integer, nullable=False, default=0)
    visible = Column(Boolean, default=False)

class Database:
    def __init__(self):
        self.DB_CONFIG = {
            "host": "",
            "user": "",
            "password": "",
            "database": "",
            "port": 5432
        }

        self.__get_config()
        self.session = self.__connect()

    def __get_config(self):
        envVars = ["POSTGRES_USER", "POSTGRES_HOST", "POSTGRES_PASSWORD", "POSTGRES_DATABASE"]

        for var in envVars:
            if not os.getenv(var):
                raise ValueError(f"Environment variable {var} not set")
            
        self.DB_CONFIG["host"] = os.getenv("POSTGRES_HOST")
        self.DB_CONFIG["user"] = os.getenv("POSTGRES_USER")
        self.DB_CONFIG["password"] = os.getenv("POSTGRES_PASSWORD")
        self.DB_CONFIG["database"] = os.getenv("POSTGRES_DATABASE")
        self.DB_CONFIG["port"] = int(os.getenv("POSTGRES_PORT")) if os.getenv("POSTGRES_PORT") and os.getenv("POSTGRES_PORT").isdigit() else 5432

    def __connect(self):
        DATABASE_URL = f"postgresql+psycopg2://{self.DB_CONFIG['user']}:{self.DB_CONFIG['password']}@{self.DB_CONFIG['host']}:{self.DB_CONFIG['port']}/{self.DB_CONFIG['database']}"

        for count in range(1, 6):
            try:
                engine = create_engine(DATABASE_URL)
                Base.metadata.create_all(engine)
                Session = sessionmaker(bind=engine)
                session = Session()
                return session
            except Exception as e:
                print(f"Error connecting to the database: {e}")
                print(f"Retrying in {count * 5} seconds...")
                time.sleep(count * 5)

        raise Exception("Could not connect to the database")
    
    def new_challenge(self, kwargs) -> None:
        category_id = self.__create_category(kwargs.get("category"))
        kwargs["category_id"] = category_id
        del kwargs["category"]

        hints = kwargs.get("hint_objects", [])
        del kwargs["hint_objects"]

        challenge = self.session.query(Challenge).filter_by(chall_name=kwargs["chall_name"]).first()
        if not challenge:
            challenge = Challenge(**kwargs)
            self.session.add(challenge)
        else:
            for key, value in kwargs.items():
                setattr(challenge, key, value) if value else None
        
        self.session.commit()
        self.__add_hints(challenge.chall_id, hints)

    def __add_hints(self, challenge_id, hints) -> None:
        for hint in hints:
            old_hint = self.session.query(Hint).filter_by(chall_id=challenge_id, hint=hint["hint"]).first()
            if old_hint:
                for key, value in hint.items():
                    setattr(old_hint, key, value) if value else None
                continue
            hint = Hint(chall_id=challenge_id, **hint)
            self.session.add(hint)

        self.session.commit()

    def __create_category(self, category_name) -> None:
        category = self.session.query(Category).filter_by(category_name=category_name).first()
        if category:
            return category.category_id
        category = Category(category_name=category_name)
        self.session.add(category)
        self.session.commit()

        return category.category_id