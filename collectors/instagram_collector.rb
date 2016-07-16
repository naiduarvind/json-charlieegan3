class InstagramCollector
  def self.collect(client_id, client_secret, token)
    Instagram.configure do |config|
      config.client_id = client_id
      config.client_secret = client_secret
    end

    client = Instagram.client(access_token: token)
    image_hash = client.user_recent_media.first.select do |key,_|
      %w(link created_time location images).include? key
    end
    return format(image_hash)
  end

  private

  def self.format(image_hash)
    image_hash['created_at'] = Time.at(image_hash['created_time'].to_i).utc
    image_hash.delete('created_time')
    image_hash['images'] = Hash[image_hash['images'].map { |k, v| [k,v['url']] }]
    image_hash['location'] = image_hash['location'] ? image_hash['location']['name'] : "someplace"
    return image_hash
  end
end
