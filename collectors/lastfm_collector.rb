class LastfmCollector
  def self.collect(key, secret, username)
    Rockstar.lastfm = { api_key: key, api_secret:  secret }
    track = Rockstar::User.new(username).recent_tracks.first
    return unless track
    {
      'name' => track.name,
      'artist' => track.artist,
      'link' => track.url,
      'image' => select_image(track.images),
      'created_at' => (track.date || Time.now).utc
    }
  end

  def self.select_image(images)
    [
      images['large'],
      images['medium'],
      images['extralarge'],
      images['small'],
    ].compact.first
  end
end
