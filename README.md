# Naive Reverend

The Naive Reverend is an HTTP service for Bayesian textual classification using a [bag of words](https://en.wikipedia.org/wiki/Bag-of-words_model) model. Despite their simplicity, bag-of-words models can perform remarkably well for tasks like spam filtering and sentiment detection. They're fast and easy to implement. And maybe above all, writing one from scratch is a good way to drill [Bayes' Theorem](https://en.wikipedia.org/wiki/Bayes%27_theorem) into your head and be forced to wrestle with some of the subtleties of floating point math on a computer.

## About the name

In addition to being a statistician and philosopher, Thomas Bayes was a Presbyterian minister. In a bag of words classifier, we make the "naive" assumption that all features, in our case, words, are conditionally independent. In other words the probability of a word occuring in a class is independent of the words around it. Of course, that isn't true, but we can still build pretty good classifiers if we let ourselves make that assumption. We can then evaluate the accuracy of the classifier using a hold out test set.

# Endpoints

### /classify
### /train

## Store backends
### Redis
### In-memory
### LevelDB

## Should I use it?

Probably not. Aside from the fact that it has almost no tests, you can likely get much more accurate classification with a [backoff language model](https://en.wikipedia.org/wiki/Katz%27s_back-off_model) using a library like [kenlm](https://kheafield.com/code/kenlm/), [berkeleylm](https://code.google.com/p/berkeleylm/), or [irstlm](http://sourceforge.net/projects/irstlm/). All of these libraries use data structures that have been highly optimized for read performance and space efficiency. But they're significantly more expensive to update and retrain than just incrementing counts in a key value store.
