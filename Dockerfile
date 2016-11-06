FROM ruby

WORKDIR /app

# App
COPY Gemfile Gemfile.lock ./
RUN gem install bundler
RUN bundle install

COPY . /app

CMD ruby status.rb
