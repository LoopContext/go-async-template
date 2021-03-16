### Configure environment
- npm install -g @asyncapi/generator

### While building template or extending it use:
copy async API file to root folder and kick the following
ag ./server_order.yaml ./ --force-write --output ./out/ --watch-template --install
(if fails to pickup changes, bump package.json version and delete node_modules)

### To deploy
* Bump version in package.json
* Set npm repository -
  npm config set registry https://integrapartners.jfrog.io/artifactory/api/npm/npm/
  
* Login to npm -
  npm login

* Enter your user/password

* Publish
  npm publish
