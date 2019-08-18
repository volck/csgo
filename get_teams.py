from html.parser import HTMLParser
from requests import get


class MyHTMLParser(HTMLParser):

    def __init__(self):
        HTMLParser.__init__(self)
        self.recording = 0
        self.data = []

    def handle_starttag(self, tag, attrs):
        if tag == 'span':
            for name, value in attrs:
                if name == 'class' and value == 'team-name':
                    self.recording = 1

    def handle_endtag(self, tag):
        if tag == 'span':
            self.recording -= 1

    def handle_data(self, data):
        if self.recording == 1:
            if len(data.strip()) > 2:
                self.data.append(data.strip())


def grab_team_names():
    p = MyHTMLParser()
    req = get(
        'https://www.gamer.no/turneringer/telialigaen-counter-strike-go-hosten-2019/6095/tabeller/')

    p.feed(req.text)
    return p.data


for team in grab_team_names():
    print(team)
