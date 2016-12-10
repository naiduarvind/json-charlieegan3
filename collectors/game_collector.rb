class GameCollector
  def self.collect(steam_username, psn_username, sc2_url)
    [
      psn_latest(psn_username),
      steam_latest(steam_username),
      sc2_latest(sc2_url),
    ].compact
  end

  def self.steam_latest(steam_username)
    doc = Nokogiri::HTML(open("https://steamcommunity.com/id/#{steam_username}"))
    date = doc.css('.game_info_details').first
    return unless date
    date_string = date.children.last.text.strip
    {
      network_icon: "https://steamcommunity.com/favicon.ico",
      action: "https://steamcommunity.com/id/#{steam_username}",
      game: doc.css('.game_name').first.text,
      time: Time.parse(date_string.gsub("last played on ", "")).ago_in_words.gsub(/ and \w+ \w+/, ''),
    }
  end

  def self.psn_latest(psn_username)
    doc = Nokogiri::HTML(open("https://psnprofiles.com/#{psn_username}"))
    thing = {
      network_icon: "https://www.playstation.com/favicon.ico",
      action: "https://my.playstation.com/#{psn_username}",
      game: doc.css('#gamesTable a.title').first.text,
      time: nil
    }
  end

  def self.sc2_latest(sc2_url)
    doc = Nokogiri::HTML(open(sc2_url))
    date = doc.css('td.align-right').first.text.strip
    ago_string = Time.strptime(date, '%d/%m/%Y').ago_in_words.gsub(/ and \w+ \w+/, '')
    {
      network_icon: "https://eu.battle.net/sc2/static/images/icons/favicon.ico",
      action: sc2_url.gsub("matches", ""),
      game: "Starcraft 2",
      time: ago_string,
    }
  end
end
