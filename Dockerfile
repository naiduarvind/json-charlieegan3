FROM ruby

WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN gem install bundler
RUN bundle install

ADD . /app

CMD ruby status.rb
