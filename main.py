import json
from requests import get

from bs4 import BeautifulSoup


def json_config():
    with open('sites.json') as json_blob:
        json_data = json.load(json_blob)

    return json_data


def main():
    data = json_config()

    for site in data:
        response = get(data[site]['url'])
        print(response.text)
        soup = BeautifulSoup(response.content, 'html.parser')
        result = soup.find(id=data[site]['div'])

        print(result)


if __name__ == '__main__':
    main()
