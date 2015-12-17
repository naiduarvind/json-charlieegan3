require 'json'
require 'open-uri'

require 'aws-sdk'
require 'aws-sdk-resources'
require 'instagram'
require 'rockstar'
require 'strava/api/v3'
require 'twitter'
require 'time_ago_in_words'

require './collectors/github_collector'
require './collectors/twitter_collector'
require './collectors/strava_collector'
require './collectors/lastfm_collector'
require './collectors/instagram_collector'

def ago_string(time)
  time.ago_in_words.gsub(/ and \w+ \w+/, '')
end

status = JSON.parse(open(ENV['STATUS_URL']).read)
status.delete('metadata')

twitter_credentials = [
  ENV['TWITTER_KEY'], ENV['TWITTER_SECRET'], ENV['TWITTER_ACCESS_TOKEN_KEY'],
  ENV['TWITTER_ACCESS_TOKEN_SECRET']
]

instagram_credentials = [
  ENV['INSTAGRAM_CLIENT_ID'], ENV['INSTAGRAM_CLIENT_SECRET'],
  ENV['INSTAGRAM_ACCESS_TOKEN']
]

lastfm_credentials = [ENV['LASTFM_KEY'], ENV['LASTFM_SECRET'], ENV['USERNAME']]

status['activity'] = StravaCollector.collect(ENV['STRAVA_TOKEN'])
status['commit'] = GitHubCollector.collect(ENV['USERNAME'])
status['image'] = InstagramCollector.collect(*instagram_credentials)
status['track'] = LastfmCollector.collect(*lastfm_credentials)
status['tweet'] = TwitterCollector.collect(ENV['USERNAME'], *twitter_credentials)

status.each do |k,_|
  status[k].merge!('created_ago' => ago_string(status[k]['created_at']))
end

status['metadata'] = { created_at: Time.new.utc }

Aws.config.update({
  region: ENV['AWS_REGION'],
  credentials: Aws::Credentials.new(ENV['AWS_KEY'], ENV['AWS_SECRET'])
})
bucket = Aws::S3::Resource.new.bucket(ENV['AWS_BUCKET'])

path = 'status.json'
obj = bucket.object(path)
obj.put(body: status.to_json, acl: 'public-read')
