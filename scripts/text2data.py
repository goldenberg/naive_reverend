import nltk
import os
import sys
import json

def main():
	directory = sys.argv[1]
	label = os.path.split(directory)[-1]

	for filename in os.listdir(directory):
		if filename.startswith('.'):
			continue
		tokens = []
		for line in open(os.path.join(directory, filename)):
			tokens.extend(tokenize(line))

		print json.dumps({
			'class': label,
			'features': tokens,
		})


PUNCTUATION = list(r',.;\'"!@#$%^&*()-_?')

def tokenize(text):
	def any_punct(token):
		token_set = token
		return any(
			p in token_set
			for p in PUNCTUATION)

	return [
		token for token in
		nltk.tokenize.wordpunct_tokenize(text)
		if not any_punct(token) and len(token) > 2
	]
if __name__ == '__main__':
	main()
