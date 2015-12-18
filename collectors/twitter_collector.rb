class TwitterCollector
  def self.collect(username, key, secret, token_key, token_secret)
    client = initialize_client(key, secret, token_key, token_secret)

    tweet = client.user_timeline(username).reject do |t|
      t.reply? || t.retweet? || t.media? || t.uris? || t.user_mentions?
    end.first

    {
      'text' => tweet.text,
      'location' => tweet.place.full_name,
      'link' => tweet.url.to_s,
      'created_at' => Time.parse(tweet.created_at.to_s).utc
    }
  end

  private

  def self.initialize_client(key, secret, token_key, token_secret)
    Twitter::REST::Client.new do |config|
      config.consumer_key        = key
      config.consumer_secret     = secret
      config.access_token        = token_key
      config.access_token_secret = token_secret
    end
  end
end
