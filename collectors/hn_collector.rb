class HackerNewsCollector
  def self.collect(user_id)
    doc = Nokogiri::HTML(open("https://news.ycombinator.com/submitted?id=#{user_id}").read)
    html = doc.css('.itemlist tr').take(2)
    link = html.first.css('a')[1]
    {
      title: link.text,
      url: link["href"],
      comments: "https://news.ycombinator.com/" + html.last.css('a').last["href"],
      created_ago: html.last.css('span')[1].text,
    }
  end
end
