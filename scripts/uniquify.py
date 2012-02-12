import json
import sys

def main():
	for line in sys.stdin:
		x = json.loads(line)
		x['features'] = list(set(x['features']))
		print json.dumps(x)

if __name__ == '__main__':
	main()