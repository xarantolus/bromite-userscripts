names="cosmetic cosmetic-lite idcac twitter"

for name in $names; do
	# Basically create "$name.html" by replacing {{URL}} with "https://github.com/xarantolus/bromite-userscripts/releases/latest/download/{{name}}.user.js"
	sed "s/{{URL}}/https:\/\/github.com\/xarantolus\/bromite-userscripts\/releases\/latest\/download\/${name}.user.js/g" template.html > "$name.html"
done

