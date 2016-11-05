while true
do
  echo "Starting..."
  ruby /app/status.rb
  echo "Completed."
  echo "Running in 10 minutes"
  sleep 600
done
