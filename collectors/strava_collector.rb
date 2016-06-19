class StravaCollector
  def self.collect(token)
    client = Strava::Api::V3::Client.new(access_token: token)

    keys = %w(id name distance moving_time location_city start_date)
    activity = client.list_athlete_activities.first.select do |key, _|
      keys.include? key
    end

    activity['created_at'] = Time.parse(activity['start_date'])
    activity.delete('start_date')
    activity['location'] = activity['location_city']
    activity.delete('location_city')
    activity['link'] = "https://www.strava.com/activities/#{activity['id']}"
    activity.delete('id')

    activity['distance'] = (activity['distance'] / 1000).round(1)

    time = activity['moving_time'] / 60.0
    mins = time.to_i
    activity['moving_time'] = "#{mins} minutes #{((time - mins) * 60).to_i} seconds"

    id = client.retrieve_current_athlete["id"]
    activity['ytd'] =
      client.totals_and_stats(id)["ytd_run_totals"]["distance"] / 1000

    return activity
  end
end
