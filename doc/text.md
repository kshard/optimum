# Nearest Neighborhood Search

## Example use-cases

The data structure clusters textual content using nearest neighborhood and
provides efficient lookup. Typical usage:

* Text Search and Retrieval: The struct can be used in search engines and document
retrieval systems to quickly find similar documents on high-dimensional
text embeddings.

* Content-Based Recommendations: The struct is useful for finding similar items in
recommendation systems, such as finding related products, movies, or music
tracks based on embeddings vectors. It helps quickly locate users or items
with similar behavior patterns.

* Personalized Content: When a system needs to recommend personalized content
(e.g., news articles, blog posts), the struct can quickly find the most relevant
content based on a user's preferences or behavior.

* Semantic Search: In NLP, The struct is used to find semantically similar phrases,
sentences, or documents by comparing embeddings generated by text  models.

* Chatbots and Conversational AI: The struct can be used to match user queries to a set
of predefined responses or intents based on similarity.


## Create data structure instance

```bash
optimum text create -u $HOST -n <name> -j path/to/config.json
```

The algorithm "text" is an approximation nearest neighbor search of natural
language content. It enhances the usability of "approximate nearest
neighbor search in high-dimensional spaces" by integrating embeddings and
supporting indexing and retrieval of textual blocks instead of pure vectors.

Config algorithm through primary parameters:
- "embeddings.model" the model id of the model to calculate embeddings vector.

- "embeddings.dimension" is a size of output embeddings vector.

- if "hnsw.*" is defined, it enables the use of the Hierarchical Navigable
  Small World (HNSW) algorithm for approximation nearest neighbor search.

- "hnsw.M" and "hnsw.M0" controls the maximum number of connections per node,
  balancing between memory usage and search efficiency.  M0 defines
  the connection density on the graph's base layer, while M regulates it on
  the intermediate layers.

- "hnsw.efConstruction" determines the number of candidate nodes evaluated
  during graph construction, influencing both the construction time and
  the accuracy of the graph.

- "hnsw.surface" is vector distance function.

Example configuration and default values:	

```json
{
  "embeddings": {
    "model": "amazon.titan-embed-text-v2:0",
    "dimension": 256,       // enum {256, 512, 1024} 
  },
  "hnsw": {
    "m":  8,                // number in range of [4, 1024]
    "m0": 64,               // number in range of [4, 1024]
    "efConstruction": 200,  // number in range of [200, 1000]
    "surface": "cosine"     // enum {"cosine", "euclidean"}
  }
}
```

Supported embeddings models:
* `amazon.titan-embed-text-v2:0`


## Writing to data structure instance (batch mode)

The batch writing consist of two phases - data upload followed by a commit.

```bash
# Upload data into server.
optimum text upload -u $HOST -n <name> path/to/data.txt

# Commit uploaded data, making it available online.
optimum text commit -u $HOST -n <name>
```

The upload supports two file formats: text or json. The type of the file is
determined by the file extension.

**Text files (.txt)**

Each line of the file is a text block to be indexed as whole. The format does
not carry on any metadata, it is tailored for pure text processing:

```
his garret was under the roof of a high, ... cupboard than a room.
it's simply a fantasy to amuse myself; a plaything!
```

**Json files (.json)**

Each line of the file is a json object that carries on the text and metadata to
be indexed:

```json
{"text": "his garret was...", "isPartOf": "https://example.com/abc123", ...}
{"text": "a fantasy to amuse...", "isPartOf": "https://example.com/abc456", ...}
```

The full schema of json object is following:

```json
{
  "text": "...",       // short text block less than 4KB.
  "isPartOf": "...",   // URL of the original doc from which the text is derived.
  "headline": ["..."], // headline(s) of the text.
  "keywords": ["..."], // relevant keywords for the text.
  "links": ["..."]     // external URIs associated with the text. 
}
```