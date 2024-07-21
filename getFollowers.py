#!/usr/bin/python3
import logging
logging.info("Hello")

import re
import csv
from time import sleep
import os
import sys
import pathlib
from timeit import default_timer as timer
import datetime
import os
from dotenv import load_dotenv

load_dotenv()

import urllib3
import instaloader

username = os.getenv('INSTA_USERNAME')
password = os.getenv('INSTA_PASSWORD')

# Get instance
L = instaloader.Instaloader()

# Login or load session

L.login(username, password)        # (login)


http = urllib3.PoolManager()

start = timer()
curr = str(datetime.datetime.now())   
try: os.remove("scraped.txt")
except: pass

f = open('input.txt','r')
accounts = f.read()
PROFILE = accounts.split('\n')

for ind in range(len(PROFILE)):
    pro = PROFILE[ind]
    try:
        print('\n\nGetting followers from',pro)
        filename = 'scraped.txt'
        
        profile = instaloader.Profile.from_username(L.context, pro)
        main_followers = profile.followers
        count = 0
        total = 0
        # Print list of followees
        for num, person in enumerate(profile.get_followers()):
            try:
                total+=1
                username = person.username
                try:
                    follower_profile = instaloader.Profile.from_username(L.context, username)
                except Exception as e:
                    print(e)

                print('Username:',username)
                logging.info('Username:',username)
                with open(filename,'a',newline='') as csvf:
                    csvf.write(username+'\n')

                print('--------------------------------------------------------------------------------\nTotal followers scraped:',total,' out of',main_followers)
                logging.info('--------------------------------------------------------------------------------\nTotal followers scraped:',total,' out of',main_followers)

                print('Time:',str(datetime.timedelta(seconds=(timer()-start))))
                logging.info('Time:',str(datetime.timedelta(seconds=(timer()-start))))
                
            except Exception as e:
                print(e)
                logging.info(e)

            if num == int(sys.argv[1])-1: break
    except Exception as e:
        print('Skipping',pro)
        print(e)
        logging.info('Skipping',pro)