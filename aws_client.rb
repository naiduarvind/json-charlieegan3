class AwsClient
  def initialize(key, secret, region, distribution)
    @key, @secret, @region, @distribution = key, secret, region, distribution
    Aws.config.update({
      region: region,
      credentials: Aws::Credentials.new(key, secret)
    })
  end

  def post(bucket, file, body)
    client = Aws::S3::Client.new(
      region: @region,
      credentials: Aws::Credentials.new(@key, @secret),
    )
    s3 = Aws::S3::Resource.new(client: client)
    bucket = s3.bucket(ENV["AWS_BUCKET"])
    bucket.put_object(
      key: file,
      body: body,
      acl: "public-read",
      content_type: "application/json",
      expires: (Time.new + 60*10).httpdate,
      cache_control: "public, max-age=600"
    )
  end

  def invalidate(path)
    cloudfront = Aws::CloudFront::Client.new(
      region: @region,
      credentials: Aws::Credentials.new(@key, @secret),
    )

    cloudfront.create_invalidation({
      distribution_id: @distribution,
      invalidation_batch: {
        paths: {
          quantity: 1,
          items: [path],
        },
        caller_reference: "json-charlieegan3",
      },
    })
  end
end

