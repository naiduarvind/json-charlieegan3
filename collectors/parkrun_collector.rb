require 'pry'
class ParkrunCollector
  def self.collect(barcode)
    # only check 1x per week
    return nil unless Time.now.wday == 0 && Time.now.hour == 12

    url = "https://www.parkrun.org.uk/results/athleteresultshistory/?athleteNumber=#{barcode}"
    doc = Nokogiri::HTML(open(url, "User-Agent" => "Chrome").read)
    latest_row = doc.at_css("#results tbody tr")

    link = latest_row.children.first.css("a").last
    {
      "location" => link.text.strip.scan(/^\w+/).first,
      "link" => link["href"],
      "time" => latest_row.children[4].text,
      "created_at" => Time.strptime(latest_row.children[1].text, '%d/%m/%Y'),
    }
  end
end
