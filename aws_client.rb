class AwsClient
  def initialize(key, secret, region)
    @key, @secret, @region = key, secret, region
    Aws.config.update({
      region: region,
      credentials: Aws::Credentials.new(key, secret)
    })
  end

  def post(bucket, file, body)
    bucket = Aws::S3::Resource.new.bucket(bucket)

    obj = bucket.object(file)
    obj.put(body: body, acl: 'public-read')
  end
end

