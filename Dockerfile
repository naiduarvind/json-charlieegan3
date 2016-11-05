FROM ruby

WORKDIR /app

# App
COPY Gemfile Gemfile.lock ./
RUN gem install bundler
RUN bundle install

COPY . /app

CMD ["bash", "-e", "entrypoint.sh"]
