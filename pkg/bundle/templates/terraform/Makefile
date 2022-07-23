.PHONY: plan

 plan: ## Plans all examples
 	mass bundle build
 	ruby ./tf_helper.rb plan | bash

 apply: ## Applies all examples
 	mass bundle build
 	ruby ./tf_helper.rb apply | bash

 destroy: ## Destroys all examples
 	mass bundle build
 	ruby ./tf_helper.rb destroy | bash
