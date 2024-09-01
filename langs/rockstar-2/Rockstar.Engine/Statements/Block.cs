using System.Text;
using Rockstar.Engine.Expressions;

namespace Rockstar.Engine.Statements;

public class Block {
	public List<Statement> Statements { get; } = [];
	public Block() { }
	public Block Concat(Block tail) {
		this.Statements.AddRange(tail.Statements);
		return this;
	}
	public Block(params Statement[] statements) => this.Statements.AddRange(statements);

	public override string ToString() => Print(new()).ToString();

	public StringBuilder PrintTopLevel(StringBuilder sb) {
		foreach (var stmt in Statements) stmt.Print(sb, "");
		return sb;
	}

	public StringBuilder Print(StringBuilder sb, string prefix = "") {
//		if (Statements.Count == 1) return Statements[0].Print(sb, (prefix == "" ? "" : prefix + Expression.INDENT));
//		if (prefix != "") {
//			sb.Append(prefix).AppendLine("block:");
			foreach (var stmt in Statements) stmt.Print(sb, prefix + "│ ");
			return sb.Append(prefix).AppendLine("└──────────");
//		}
		foreach (var stmt in Statements) stmt.Print(sb, prefix == "" ? "" : prefix + Expression.INDENT);
		return sb;
	}
}