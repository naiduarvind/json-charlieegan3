class LetterboxdCollector
  def self.collect(username)
    diary_url = "https://letterboxd.com/#{username}/films/diary/"
    doc = Nokogiri::HTML(open(diary_url).read)
    row = doc.at_css("#diary-table .diary-entry-row")

    slug = row.at_css('.td-film-details div')['data-film-slug']
    cover_doc = Nokogiri::HTML(open("https://letterboxd.com#{slug}image-150/"))

    {
      'title' => row.at_css('.td-film-details').text.strip,
      'cover' => cover_doc.at_css('img')['src'],
      'link' => diary_url,
      'created_at' => Time.parse(row.children[1..3].map(&:text).join.strip),
    }
  end
end
