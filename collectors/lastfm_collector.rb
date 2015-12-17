class LastfmCollector
  def self.collect(key, secret, username)
    Rockstar.lastfm = { api_key: key, api_secret:  secret }
    track = Rockstar::User.new(username).recent_tracks.first
    {
      name: track.name,
      artist: track.artist,
      link: track.url,
      images: track.images
    }
  end
end
