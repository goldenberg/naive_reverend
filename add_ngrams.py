import itertools
import json
import sys

def main():
	for line in sys.stdin:
		x = json.loads(line)
		bigrams = (
			' '.join(tokens) for tokens in
			itertools.izip(x['features'], x['features'][1:])
		)

		trigrams = (
			' '.join(tokens) for tokens in
			itertools.izip(x['features'], x['features'][1:], x['features'][2:])
		)

		quadgrams = (
			' '.join(tokens) for tokens in
			itertools.izip(x['features'], x['features'][1:], x['features'][2:], x['features'][2:])
		)

		x['features'].extend(bigrams)
		x['features'].extend(trigrams)
		x['features'].extend(quadgrams)
		print json.dumps(x)


if __name__ == '__main__':
	main()