using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Text.Json.Nodes;

var _removeSet = new HashSet<string>();
var _removeCount = 0;

var listPath = "trimmer.txt";
var lines = File.ReadAllLines(listPath);

foreach (var filename in lines)
{
	var path = $"/out/{filename}";
	if (!File.Exists(path))
	{
		Console.Error.WriteLine($"File specified in {listPath} not found: {path}");
	}

	File.Delete(path);
	_removeSet.Add(filename);
}

var jsonPath = "/out/f-sharp.deps.json";
var text = File.ReadAllText(jsonPath);
var data = JsonNode.Parse(text);
Walk(data);
File.WriteAllText(jsonPath, data.ToJsonString(new JsonSerializerOptions { WriteIndented = true }));
Console.WriteLine($"Based on {listPath}, deleted {_removeSet.Count} files and removed {_removeCount} entries in {jsonPath}.");

void Walk(JsonNode node)
{
	switch (node)
	{
		case JsonObject obj:
			WalkObject(obj);
			break;
		case JsonArray array:
			foreach (var child in array)
			{
				Walk(child);
			}
			break;
	}
}

void WalkObject(JsonObject obj)
{
	var toRemove = new List<string>();
	foreach (var property in obj)
	{
		if (_removeSet.Contains(property.Key))
		{
			toRemove.Add(property.Key);
			_removeCount++;
			continue;
		}

		Walk(property.Value);
	}

	foreach (var key in toRemove)
	{
		obj.Remove(key);
	}
}
