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

status = JSON.parse(open(ENV['STATUS_URL']).read)

twitter_credentials = {
  key: ENV['TWITTER_KEY'],
  secret: ENV['TWITTER_SECRET'],
  token_key: ENV['TWITTER_ACCESS_TOKEN_KEY'],
  token_secret: ENV['TWITTER_ACCESS_TOKEN_SECRET']
}

status['activity'] = StravaCollector.collect(ENV['STRAVA_TOKEN'])
status['commit'] = GitHubCollector.collect(ENV['USERNAME'])
status['image'] = InstagramCollector.collect(ENV['INSTAGRAM_CLIENT_ID'], ENV['INSTAGRAM_CLIENT_SECRET'], ENV['INSTAGRAM_ACCESS_TOKEN'])
status['track'] = LastfmCollector.collect(ENV['LASTFM_KEY'], ENV['LASTFM_SECRET'], ENV['USERNAME'])
status['tweet'] = TwitterCollector.collect(ENV['USERNAME'], twitter_credentials)

status['metadata'] = { collected_date: Time.new.utc }


Aws.config.update({
  region: ENV['AWS_REGION'],
  credentials: Aws::Credentials.new(ENV['AWS_KEY'], ENV['AWS_SECRET'])
})
bucket = Aws::S3::Resource.new.bucket(ENV['AWS_BUCKET'])

path = 'status.json'
obj = bucket.object(path)
obj.put(body: status.to_json, acl: 'public-read')
