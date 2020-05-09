using System;
using System.Collections.Generic;
using System.IO;
using Newtonsoft.Json.Linq;

class Trimmer
{
	static HashSet<string> _removeSet = new HashSet<string>();
	static int _removeCount;

	static void Main(string[] args)
	{
		var listPath = "ExtraTrimmingList.txt";
		var lines = File.ReadAllLines(listPath);

		foreach (var filename in lines)
		{
			var path = $"/compiler/{filename}";
			if (!File.Exists(path))
			{
				Console.Error.WriteLine($"File specified in {listPath} not found: {path}");
			}

			File.Delete(path);
			_removeSet.Add(filename);
		}

		var jsonPath = "/compiler/Compiler.deps.json";
		var text = File.ReadAllText(jsonPath);
		var data = JObject.Parse(text);
		Walk(data);
		File.WriteAllText(jsonPath, data.ToString());
		Console.WriteLine($"Based on {listPath}, deleted {_removeSet.Count} files and removed {_removeCount} entries in {jsonPath}.");
	}

	static void Walk(JToken token)
	{
		switch (token)
		{
			case JObject obj:
				WalkObject(obj);
				break;
			case JArray array:
				foreach (var child in array.Children())
				{
					Walk(child);
				}
				break;
			case JProperty property:
				Walk(property.Value);
				break;
		}
	}

	static void WalkObject(JObject obj)
	{
		var toRemove = new List<JProperty>();
		foreach (var property in obj.Properties())
		{
			if (_removeSet.Contains(property.Name))
			{
				toRemove.Add(property);
				_removeCount++;
				continue;
			}

			Walk(property);
		}

		foreach (var property in toRemove)
		{
			property.Remove();
		}
	}
}
