using System.Text;
using Rockstar.Engine.Expressions;
using Rockstar.Engine.Values;
using Array = Rockstar.Engine.Values.Array;

namespace Rockstar.Engine.Statements;

public class Output(Expression expr, string suffix = "") : ExpressionStatement(expr) {
	public string Suffix { get; } = suffix;

	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.Append(prefix).AppendLine("output: ");
		return Expression.Print(sb, prefix + INDENT);
	}
}

public class Debug(Expression expr) : ExpressionStatement(expr) {
	public override StringBuilder Print(StringBuilder sb, string prefix) {
		sb.Append(prefix).Append("debug: ");
		return Expression switch {
			Lookup lookup => sb.AppendLine(lookup.ToString()),
			Value value => sb.AppendLine(value.ToString()),
			_ => Expression.Print(sb.AppendLine(), prefix + INDENT)
		};
	}

}