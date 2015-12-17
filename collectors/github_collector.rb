class GitHubCollector
  def self.collect(username)
    feed = fetch_feed_for_user(username)
    feed = filter_feed_for_suitable_commits(feed)
    feed.map { |c| format_commit(c) }.first
  end

  private

  def self.fetch_feed_for_user(username)
    JSON.parse(open("https://api.github.com/users/#{username}/events").read)
  end

  def self.filter_feed_for_suitable_commits(feed)
    feed.select! { |e| e['type'] == 'PushEvent' && e['payload']['commits']}
    feed.map! do |e|
      e['payload']['commits'].map do |c|
        c.merge!('created_at' => e["created_at"])
      end
    end
    feed.flatten!
    feed.reject { |e| e["message"].include? "Automated commit"}
  end

  def self.format_commit(commit)
    commit.select! { |k,_| %w(message url created_at).include? k }
    commit['url'].gsub!(/api\.|repos\//, '')
    commit['url'].gsub!('commits', 'commit')
    commit['created_at'] = Time.parse(commit['created_at'])
    return commit
  end
end
