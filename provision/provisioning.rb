node["recipes"] = node["recipes"] || []

node["recipes"].each do |recipe|
  include_recipe recipe
end
