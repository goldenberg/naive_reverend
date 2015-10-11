# Naive Reverend

The Naive Reverend is an HTTP service for Bayesian textual classification using a [bag of words](https://en.wikipedia.org/wiki/Bag-of-words_model) model. Despite their simplicity, bag-of-words models can perform remarkably well for tasks like spam filtering and sentiment detection. They're fast and easy to implement. And maybe above all, writing one from scratch is a good way to drill [Bayes' Theorem](https://en.wikipedia.org/wiki/Bayes%27_theorem) into your head and be forced to wrestle with some of the subtleties of floating point math on a computer.

## About the name

In addition to being a statistician and philosopher, Thomas Bayes was a Presbyterian minister.

# Endpoints

### /classify
### /train

## Store backends
### Redis
### In-memory
### LevelDB

## Should I use it?

Probably not. Aside from the fact that it has almost no tests, you can likely get much more accurate classification with similar read performance characteristics by using a [backoff language model](https://en.wikipedia.org/wiki/Katz%27s_back-off_model) using a library like [kenlm](https://kheafield.com/code/kenlm/), [berkeleylm](https://code.google.com/p/berkeleylm/), or [irstlm](http://sourceforge.net/projects/irstlm/). But they'll be significantly more expensive to retrain than just incrementing counts in a key value store.
