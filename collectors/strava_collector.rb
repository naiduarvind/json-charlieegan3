class StravaCollector
  def self.collect(token)
    client = Strava::Api::V3::Client.new(access_token: token)

    keys = %w(name distance moving_time start_latlng location_city start_date)
    activity = client.list_athlete_activities.first.select do |key,_|
      keys.include? key
    end
    activity['created_at'] = Time.parse(activity['start_date'])
    activity.delete('start_date')
    return activity
  end
end
