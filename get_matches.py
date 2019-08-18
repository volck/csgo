from html.parser import HTMLParser
from requests import get


class MyHTMLParser(HTMLParser):

    def __init__(self):
        HTMLParser.__init__(self)
        self.recording = 0
        self.data = []

    def handle_starttag(self, tag, attrs):
        if tag == 'ul':
            #        print("ul found", attrs)
            for name, value in attrs:
                if name == 'class' and 'no-list-style matches-round-1' in value:
                    #             print(value)
                    self.recording = 1

    def handle_endtag(self, tag):
        if tag == 'ul':
            self.recording -= 1

    def handle_data(self, data):
        if self.recording == 1:
            if len(data.strip()) <= 0 or data.strip == '-' or data.strip == 'vs':
                pass
            else:
                #            print(data.strip())
                self.data.append(data.strip())


class Match:
    def __init__(self, time, home, away):
        self.time = time
        self.home = home
        self.away = away


def grab_matches():
    p = MyHTMLParser()
    req = get(
        'https://www.gamer.no/turneringer/telialigaen-counter-strike-go-hosten-2019/6095/kamper/')

    p.feed(req.text)
    return p.data


obj = {'time': '', 'home': '', 'away': ''}
matches = []
m = grab_matches()

for index in range(len(m)):
    if index % 5 == 0:
        obj['time'] = m[index]
    if index % 5 == 1:
        obj['home'] = m[index]
        print(obj)
    if index % 5 == 3:
        obj['away'] = m[index]
        p = Match(obj['time'], obj['home'], obj['away'])
        matches.append(p)

for x in matches:
    print("[{}] {} vs {}  [*]".format(x.time, x.home, x.away))
