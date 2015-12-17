class TwitterCollector
  def self.collect(username, credentials)
    client = initialize_client(credentials)

    tweet = client.user_timeline(username).reject do |t|
      t.reply? || t.retweet? || t.media? || t.uris? || t.user_mentions?
    end.first

    {
      text: tweet.text,
      location: tweet.place.full_name,
      created_at: Time.parse(tweet.created_at.to_s).utc
    }
  end

  private

  def self.initialize_client(credentials)
    Twitter::REST::Client.new do |config|
      config.consumer_key        = credentials[:key]
      config.consumer_secret     = credentials[:secret]
      config.access_token        = credentials[:token_key]
      config.access_token_secret = credentials[:token_secret]
    end
  end
end
