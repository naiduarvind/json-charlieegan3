class InstagramCollector
  def self.collect(client_id, client_secret, token)
    Instagram.configure do |config|
      config.client_id = client_id
      config.client_secret = client_secret
    end

    client = Instagram.client(access_token: token)
    image = client.user_recent_media.first
    image_hash = image.select do |key,_|
      %w(link created_time location images).include? key
    end
  end
end
