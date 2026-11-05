# Dhoondo

Dhoondo is a website indexer cli that allows you to input a website, index the website and later on search through that website.


Currently the architecture of this project is going to be very simple:
- You will enter the website name(s)
- The main process adds the website name(s) to the rabbitmq q for background workers to pickup
- The workers pick up the message and start crawling the page
- The workers reverse index the page
- Further discovered links are added to the rabbitmq q again
- This goes on until the links are exhausted
