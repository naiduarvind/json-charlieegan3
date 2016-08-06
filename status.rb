require 'json'
require 'open-uri'

require 'aws-sdk'
require 'aws-sdk-resources'
require 'instagram'
require 'rockstar'
require 'strava/api/v3'
require 'twitter'
require 'nokogiri'
require 'time_ago_in_words'
require 'rollbar'

require './collectors/github_collector'
require './collectors/twitter_collector'
require './collectors/strava_collector'
require './collectors/lastfm_collector'
require './collectors/instagram_collector'
require './collectors/game_collector'
require './collectors/parkrun_collector'
require './collectors/hn_collector'
require './collectors/letterboxd_collector'

require './aws_client'

def ago_string(time)
  time.ago_in_words.gsub(/ and \w+ \w+/, '')
end

def most_recent_location(status)
  status.map { |_,v| [v['created_at'], v['location']] if !(v.class == Array) && v['location'] }.
    compact.
    sort_by(&:first).
    last.last
end

Rollbar.configure do |config|
    config.access_token = ENV['ROLLBAR_TOKEN']
end

begin
  status = JSON.parse(open(ENV['STATUS_URL']).read)
  status = Hash[status.map { |k, v|
    (v.class == Hash && v['created_at']) ? [k, v.merge('created_at' => Time.parse(v['created_at']))] : [k, v]
  }]
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
  status['track'] = (present = LastfmCollector.collect(*lastfm_credentials)) if present
  status['tweet'] = TwitterCollector.collect(ENV['USERNAME'], *twitter_credentials)
  status['games'] = GameCollector.collect(ENV['STEAM_USER'], ENV['PSN_USER'], ENV['SC2_URL'])
  status['parkrun'] = (parkrun_data = ParkrunCollector.collect(ENV['PARKRUN_BARCODE'])) ? parkrun_data : status['parkrun']
  status['hacker_news'] = HackerNewsCollector.collect(ENV['HACKER_NEWS_ID'])
  status['film'] = LetterboxdCollector.collect(ENV['LETTERBOXD_USERNAME'])

  status = Hash[status.map { |k, v|
    (v.class == Hash && v['created_at']) ? [k, v.merge('created_ago' => ago_string(v['created_at']))] : [k, v]
  }]

  status['metadata'] = {
    created_at: Time.new.utc,
    most_recent_location: most_recent_location(status)
  }

  client = AwsClient.new(ENV['AWS_KEY'], ENV['AWS_SECRET'], ENV['AWS_REGION'])
  client.post(ENV['AWS_BUCKET'], 'status.json', status.to_json)
rescue Exception => e
  unless e.inspect.match(/ReadTimeout|503|Over capacity|ServerError|buffer error|Server Error|timed out/)
    Rollbar.error(e)
  end
  raise e
end
