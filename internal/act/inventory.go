package act

type HostKey struct {
	Host     string `yaml:"host"`
	Key      string `yaml:"key"`
	NewValue string `yaml:"newvalue"`
}

type HostKeyList []HostKey

// RemoveKeys removes the specified keys from the data structure
// It recursively traverses the data structure and removes the keys from all maps
// It removes the line.
func RemoveKeys(data interface{}, keysToRemove []string) {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		for _, key := range keysToRemove {
			delete(v, key)
		}
		for _, value := range v {
			RemoveKeys(value, keysToRemove)
		}
	case []interface{}:
		for _, item := range v {
			RemoveKeys(item, keysToRemove)
		}
	}
}

// RemoveNulls removes all null values from the data structure
func RemoveNulls(data interface{}) {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		for key, value := range v {
			if value == nil {
				// Replace nil with an empty map to keep the key with a colon but no 'null' after
				v[key] = map[interface{}]interface{}{}
			} else {
				RemoveNulls(value)
			}
		}
	case []interface{}:
		for _, item := range v {
			RemoveNulls(item)
		}
	}
}

func UpdateHostKey(data interface{}, hostName string, key string, newValue interface{}) {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		// Check if this map has a "hosts" key
		if hosts, ok := v["hosts"]; ok {
			if hostsMap, ok := hosts.(map[interface{}]interface{}); ok {
				if hostData, ok := hostsMap[hostName]; ok {
					// Update the key-value pair for the specified host
					if hostAttrs, ok := hostData.(map[interface{}]interface{}); ok {
						hostAttrs[key] = newValue
					} else {
						// Initialize the host attributes map if it's nil or not a map
						hostAttrs = map[interface{}]interface{}{
							key: newValue,
						}
						hostsMap[hostName] = hostAttrs
					}
				}
			}
		}
		// Recursively process nested structures
		for _, value := range v {
			UpdateHostKey(value, hostName, key, newValue)
		}
	case []interface{}:
		for _, item := range v {
			UpdateHostKey(item, hostName, key, newValue)
		}
	}
}
