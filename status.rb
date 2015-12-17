require 'json'
require 'open-uri'

require 'aws-sdk'
require 'aws-sdk-resources'
require 'instagram'
require 'rockstar'
require 'strava/api/v3'
require 'twitter'

require './collectors/github_collector'
require './collectors/twitter_collector'
require './collectors/strava_collector'
require './collectors/lastfm_collector'
require './collectors/instagram_collector'

twitter_credentials = {
  key: ENV['TWITTER_KEY'],
  secret: ENV['TWITTER_SECRET'],
  token_key: ENV['TWITTER_ACCESS_TOKEN_KEY'],
  token_secret: ENV['TWITTER_ACCESS_TOKEN_SECRET']
}

data = {
  activity: StravaCollector.collect(ENV['STRAVA_TOKEN']),
  commit: GitHubCollector.collect(ENV['USERNAME']),
  image: InstagramCollector.collect(ENV['INSTAGRAM_CLIENT_ID'], ENV['INSTAGRAM_CLIENT_SECRET'], ENV['INSTAGRAM_ACCESS_TOKEN']),
  track: LastfmCollector.collect(ENV['LASTFM_KEY'], ENV['LASTFM_SECRET'], ENV['USERNAME']),
  tweet: TwitterCollector.collect(ENV['USERNAME'], twitter_credentials)
}

data.merge!({ metadata: { collected_date: Time.new } })

Aws.config.update({
  region: ENV['AWS_REGION'],
  credentials: Aws::Credentials.new(ENV['AWS_KEY'], ENV['AWS_SECRET'])
})
bucket = Aws::S3::Resource.new.bucket(ENV['AWS_BUCKET'])

path = 'status.json'
obj = bucket.object(path)
obj.put(body: data.to_json, acl: 'public-read')
