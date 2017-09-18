# grass-graph

![grass-graph-logo](grass-graph-logo.png)

Generate a PNG image of the specified GitHub account's public contribution graph.

## Use as an External Service

[Grass-Graph / Imaging your Publick Contributions](https://grass-graph.shitemil.works/)

## Try grass-graph yourself on heroku


### Setup
1. Click [![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


### Usage

Generate and Get PNG image of `github-id`.

`https://{your-heroku-app}.herokuapp.com/graph/{github-id}`

![grass-graph-normal](https://cloud.githubusercontent.com/assets/1097533/12216115/e06f1aaa-b71a-11e5-9511-18cda413027c.png)


You can specify the angle to rotate the image.

`https://{your-heroku-app}.herokuapp.com/graph/{github-id}?rotate=270`

![grass-graph-rotate](https://cloud.githubusercontent.com/assets/1097533/12216118/fb1cdb26-b71a-11e5-99f9-194185a6bcc6.png)

You can specify the size to resize the image.

`https://{your-heroku-app}.herokuapp.com/graph/{github-id}?width=350&height=50`

![grass-graph-resize](https://cloud.githubusercontent.com/assets/1097533/12216121/0a54626c-b71b-11e5-9713-d1aa6c312d0b.png)

You can also specify these two options at the same time .

`https://{your-heroku-app}.herokuapp.com/graph/{github-id}?rotate=270&width=350&height=50`

![grass-graph-both](https://cloud.githubusercontent.com/assets/1097533/12216122/178d62da-b71b-11e5-9a28-250a1a4eec76.png)

## More

In Japanese.

[GitHub の草状況を PNG 画像で返す heroku app をつくってみた - えいのうにっき](http://blog.a-know.me/entry/2016/01/09/222210)
